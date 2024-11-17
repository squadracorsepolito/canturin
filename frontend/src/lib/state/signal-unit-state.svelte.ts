import { SignalUnitService, type SignalUnit } from '$lib/api/canturin';
import { AsyncState } from './async-state.svelte';

async function loadSignalUnit(state: SignalUnitState, entityId: string) {
	state.isLoading = true;

	try {
		const sigUnit = await SignalUnitService.Get(entityId);
		state.signalUnit = sigUnit;
	} catch (error) {
		state.signalUnit = undefined;
		console.error(error);
	}

	state.isLoading = false;
}

// posso estendere non asyncstate ma l'altro!!! quindi posso implementarli uguali
class SignalUnitState extends AsyncState {
	signalUnit = $state<SignalUnit | undefined>();

	reload(entityId: string) {
		loadSignalUnit(this, entityId);
	}

	async updateName(name: string) {
		// Check if signal unit is defined
        if (!this.signalUnit) return;

		// Update the name
        try {
            const updatedSignalUnit = await SignalUnitService.UpdateName(this.signalUnit.entityId, name);
            this.signalUnit = updatedSignalUnit;
        } catch (error) {
            console.error("Error updating SignalUnit name:", error);
        }
    }

	async updateDesc(description: string) {
		// Check if signal unit is defined
		if (!this.signalUnit) return;
	
		// Update the description
		try {
			const updatedSignalUnit = await SignalUnitService.UpdateDesc(this.signalUnit.entityId, description);
			this.signalUnit = updatedSignalUnit;
		} catch (error) {
			console.error("Error updating SignalUnit description:", error);
		}
	}
	
}

export function useSignalUnit(entityId: string) {
	const state = new SignalUnitState();

	loadSignalUnit(state, entityId);

	return state;
}
