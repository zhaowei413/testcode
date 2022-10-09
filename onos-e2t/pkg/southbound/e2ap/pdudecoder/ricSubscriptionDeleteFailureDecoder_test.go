// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package pdudecoder

import (
	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2ap_commondatatypes "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-commondatatypes"
	e2ap_ies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/pdubuilder"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"gotest.tools/assert"
	"testing"
)

func Test_DecodeRicSubscriptionDeleteFailurePdu(t *testing.T) {
	procCode := v2.ProcedureCodeIDRICsubscription
	criticality := e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE
	ftg := e2ap_commondatatypes.TriggeringMessage_TRIGGERING_MESSAGE_UNSUCCESSFUL_OUTCOME
	e2apPdu, err := pdubuilder.CreateRicSubscriptionDeleteFailureE2apPdu(&types.RicRequest{
		RequestorID: 22,
		InstanceID:  6,
	}, 9,
		&e2ap_ies.Cause{
			Cause: &e2ap_ies.Cause_Transport{
				Transport: e2ap_ies.CauseTransport_CAUSE_TRANSPORT_TRANSPORT_RESOURCE_UNAVAILABLE,
			},
		})
	assert.NilError(t, err)
	assert.Assert(t, e2apPdu != nil)

	e2apPdu.GetUnsuccessfulOutcome().GetValue().GetRicSubscriptionDelete().
		SetCriticalityDiagnostics(&procCode, &criticality, &ftg,
			&types.RicRequest{
				RequestorID: 10,
				InstanceID:  20,
			}, []*types.CritDiag{
				{
					TypeOfError:   e2ap_ies.TypeOfError_TYPE_OF_ERROR_MISSING,
					IECriticality: e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE,
					IEId:          v2.ProtocolIeIDRicsubscriptionDetails,
				},
			})

	rfID, rrID, cause, pr, crit, tm, cdrID, diags, err := DecodeRicSubscriptionDeleteFailurePdu(e2apPdu)
	assert.NilError(t, err)
	assert.Equal(t, 9, int(*rfID))

	assert.Equal(t, cause.GetTransport(), e2ap_ies.CauseTransport_CAUSE_TRANSPORT_TRANSPORT_RESOURCE_UNAVAILABLE)
	assert.Equal(t, int32(*pr), int32(8))
	assert.Equal(t, int32(*crit), int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE))
	assert.Equal(t, int32(*tm), int32(e2ap_commondatatypes.TriggeringMessage_TRIGGERING_MESSAGE_UNSUCCESSFUL_OUTCOME))
	assert.Equal(t, int32(cdrID.InstanceID), int32(20))
	assert.Equal(t, int32(cdrID.RequestorID), int32(10))
	assert.Equal(t, int32(diags[0].IEId), int32(30))
	assert.Equal(t, int32(diags[0].IECriticality), int32(e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE))
	assert.Equal(t, int32(diags[0].TypeOfError), int32(e2ap_ies.TypeOfError_TYPE_OF_ERROR_MISSING))

	assert.Equal(t, 22, int(rrID.RequestorID))
	assert.Equal(t, 6, int(rrID.InstanceID))
}
