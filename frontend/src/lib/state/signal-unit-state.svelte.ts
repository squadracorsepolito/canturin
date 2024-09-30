import { SignalUnitService, type SignalUnit } from '$lib/api/canturin';
import { AsyncState } from './async-state.svelte';

async function loadSignalUnit(state: SignalUnitState, entityId: string) {
	state.isLoading = true;

	try {
		const sigUnit = await SignalUnitService.GetOpen(entityId);
		state.signalUnit = sigUnit;
	} catch (error) {
		state.signalUnit = undefined;
		console.error(error);
	}

	state.isLoading = false;
}

class SignalUnitState extends AsyncState {
	signalUnit = $state<SignalUnit | undefined>();

	reload(entityId: string) {
		loadSignalUnit(this, entityId);
	}
	
	async updateName(newName: string) {
		if (!this.signalUnit) return;

		try {
			const updatedSignalUnit = await SignalUnitService.UpdateName(this.signalUnit.entityId, newName);
			this.signalUnit = updatedSignalUnit;
		} catch (error) {
			console.error(error);
		}
	}

	async updateDesc(newDesc: string) {
		if (!this.signalUnit) return;

		try {
			const updatedSignalUnit = await SignalUnitService.UpdateDesc(this.signalUnit.entityId, newDesc);
			this.signalUnit = updatedSignalUnit;
		} catch (error) {
			console.error(error);
		}
	}

	async getInvalidNames(): Promise<string[]> {
		if (!this.signalUnit) return [];

		try {
			const invalidNames = await SignalUnitService.GetInvalidNames(this.signalUnit.entityId);
			return invalidNames || [];
		} catch (error) {
			console.error(error);
			return [];
		}
	}
}

export function useSignalUnit(entityId: string) {
	const state = new SignalUnitState();

	loadSignalUnit(state, entityId);

	return state;
}
