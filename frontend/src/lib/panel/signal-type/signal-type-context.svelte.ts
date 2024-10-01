import type { SignalType } from '$lib/api/canturin';
import { getContext, setContext } from 'svelte';

class SignalTypeState {
	signalType = $state() as SignalType;

	constructor(initial: SignalType) {
		this.signalType = initial;
	}
}

export function createSignalTypeContext(initial: SignalType) {
	const state = new SignalTypeState(initial);
	setContext(initial.entityId, state);
	return state;
}

export function getSignalTypeContext(entityId: string) {
	return getContext(entityId) as ReturnType<typeof createSignalTypeContext>;
}
