import { SignalTypeService, type SignalType } from '$lib/api/canturin';
import { AsyncState } from './async-state.svelte';

async function loadSignalType(state: SignalTypeState, entityId: string) {
	state.isLoading = true;

	try {
		const sigType = await SignalTypeService.GetOpen(entityId);
		state.signalType = sigType;
	} catch (error) {
		state.signalType = undefined;
		console.error(error);
	}

	state.isLoading = false;
}

class SignalTypeState extends AsyncState {
	signalType = $state<SignalType | undefined>();

	reload(entityId: string) {
		loadSignalType(this, entityId);
	}
}

export function useSignalType(entityId: string) {
	const state = new SignalTypeState();

	loadSignalType(state, entityId);

	return state;
}
