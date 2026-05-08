package project_evidence

import (
	"time"
)

// reconcilePrescriptionStatuses prejde všetkých pacientov a ich predpisy
// a aktualizuje status na "expired" ak doba platnosti uplynula.
func (a *Ambulance) reconcilePrescriptionStatuses() {
	now := time.Now()

	for i := range a.Patients {
		for j := range a.Patients[i].Prescriptions {
			rx := &a.Patients[i].Prescriptions[j]

			// Nemeň status ak už bol vydaný
			if rx.Status == "dispensed" {
				continue
			}

			// Skontroluj platnosť
			if rx.ValidUntil != "" {
				validUntil, err := time.Parse("2006-01-02", rx.ValidUntil)
				if err == nil && validUntil.Before(now) {
					rx.Status = "expired"
				}
			}
		}
	}
}