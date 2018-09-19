package service

import (
	"github.com/ninjadotorg/handshake-exchange/api_error"
	"github.com/ninjadotorg/handshake-exchange/bean"
	"github.com/ninjadotorg/handshake-exchange/dao"
	"github.com/ninjadotorg/handshake-exchange/integration/solr_service"
	"github.com/shopspring/decimal"
	"strings"
)

type CashService struct {
	dao     *dao.CashDao
	miscDao *dao.MiscDao
	userDao *dao.UserDao
}

func (s CashService) GetCashStore(userId string) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}

	if cashTO.Found {
		cash = cashTO.Object.(bean.CashStore)
	} else {
		ce.NotFound = true
	}

	return
}

func (s CashService) AddCashStore(userId string, body bean.CashStore) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)

	if cashTO.Error != nil {
		ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO)
		return
	}

	var err error
	if cashTO.Found {
		ce.SetStatusKey(api_error.CashStoreExists)
	} else {
		body.UID = userId
		cash = body

		err = s.dao.AddCashStore(&cash)
		if err != nil {
			ce.SetError(api_error.UpdateDataFailed, err)
			return
		}
		ce.NotFound = false
		solr_service.UpdateObject(bean.NewSolrFromCashStore(cash))
	}

	return
}

func (s CashService) UpdateCashStore(userId string, body bean.CashStore) (cash bean.CashStore, ce SimpleContextError) {
	cashTO := s.dao.GetCashStore(userId)

	if ce.FeedDaoTransfer(api_error.GetDataFailed, cashTO) {
		return
	}
	cash = cashTO.Object.(bean.CashStore)

	err := s.dao.UpdateCashStore(&cash)
	if ce.SetError(api_error.AddDataFailed, err) {
		return
	}
	solr_service.UpdateObject(bean.NewSolrFromCashStore(cash))

	return
}

func (s CashService) GetProposeCashOrder(amountStr string, currency string, fiatCurrency string) (offer bean.CashStoreOrder, ce SimpleContextError) {
	cryptoRateTO := s.miscDao.GetCryptoRateFromCache(currency, bean.INSTANT_OFFER_PROVIDER_COINBASE)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, cryptoRateTO) {
		return
	}
	cryptoRate := cryptoRateTO.Object.(bean.CryptoRate)

	systemFeeTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_CASH_BUY_CRYPTO)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, systemFeeTO) {
		return
	}
	systemFee := systemFeeTO.Object.(bean.SystemFee)

	storeFeeTO := s.miscDao.GetSystemFeeFromCache(bean.FEE_KEY_CASH_STORE_SELL_CRYPTO)
	if ce.FeedDaoTransfer(api_error.GetDataFailed, systemFeeTO) {
		return
	}
	storeFee := storeFeeTO.Object.(bean.SystemFee)

	amount, _ := decimal.NewFromString(amountStr)
	percentage, err := CreditServiceInst.GetCreditPoolPercentageByCache(currency, amount)
	if err != nil {
		if strings.Contains(err.Error(), "not enough") {
			ce.SetStatusKey(api_error.CreditOutOfStock)
			return
		}
		ce.SetError(api_error.GetDataFailed, err)
		return
	}

	price := decimal.NewFromFloat(cryptoRate.Buy).Round(2)
	totalWOFee := amount.Mul(price)

	feePercentage := decimal.NewFromFloat(systemFee.Value).Round(10)
	storeFeePercentage := decimal.NewFromFloat(storeFee.Value).Round(10)
	externalFeePercentage := decimal.NewFromFloat(float64(percentage)).Div(decimal.NewFromFloat(100))

	total, internalFee := dao.AddFeePercentage(totalWOFee, externalFeePercentage)
	_, cashFee := dao.AddFeePercentage(totalWOFee, storeFeePercentage)
	_, externalFee := dao.AddFeePercentage(totalWOFee, externalFeePercentage)

	total = total.Add(cashFee)
	total = total.Add(externalFee)

	offer.FiatAmount = total.Round(2).String()
	offer.FiatCurrency = bean.USD.Code

	offer.Amount = amountStr
	offer.Currency = currency
	offer.Price = price.String()
	offer.Fee = internalFee.Round(2).String()
	offer.FeePercentage = feePercentage.String()
	offer.StoreFee = cashFee.Round(2).String()
	offer.StoreFeePercentage = storeFeePercentage.String()
	offer.ExternalFee = externalFee.Round(2).String()
	offer.ExternalFeePercentage = externalFeePercentage.String()

	return
}
