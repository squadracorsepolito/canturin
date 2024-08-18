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
}

export function useSignalUnit(entityId: string) {
	const state = new SignalUnitState();

	loadSignalUnit(state, entityId);

	return state;
}
