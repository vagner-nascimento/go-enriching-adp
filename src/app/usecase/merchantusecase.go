package appusecase

import (
	appentity "github.com/vagner-nascimento/go-enriching-adp/src/app/entity"
	"github.com/vagner-nascimento/go-enriching-adp/src/app/interface"
	"github.com/vagner-nascimento/go-enriching-adp/src/channel"
)

func getMerchantEnrichmentData(acc appentity.Account, repo appinterface.AccountDataHandler) <-chan interface{} {
	accCh := make(chan interface{})
	affCh := make(chan interface{})

	go func() {
		defer close(accCh)

		var (
			err  error
			mAcc []appentity.MerchantAccount
		)

		if mAcc, err = repo.GetMerchantAccounts(acc.Id); err != nil {
			accCh <- err
			return
		}

		accCh <- mAcc
	}()

	go func() {
		defer close(affCh)

		var (
			err error
			aff appentity.Affiliation
		)

		if aff, err = repo.GetAffiliation(acc.Id); err != nil {
			affCh <- err
			return
		}

		affCh <- aff
	}()

	return channel.Multiplex(accCh, affCh)
}

func enrichMerchantAccount(acc *appentity.Account, mAccounts []appentity.MerchantAccount, aff *appentity.Affiliation) {
	if mAccounts != nil {
		for _, merAcc := range mAccounts {
			acc.AddMerchantAccount(merAcc)
		}
	}

	if aff != nil {
		acc.LegalDocument = &aff.LegalDocument
	}
}
