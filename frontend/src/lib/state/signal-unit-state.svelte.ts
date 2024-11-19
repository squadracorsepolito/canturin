import { SignalUnitService, type SignalUnit } from '$lib/api/canturin';
import { HistorySignalTypeModify } from '$lib/api/events';
import { EntityState } from './entity-state.svelte';
import { StateProvider } from './state-provider.svelte';

const provider = new StateProvider(
	(signalUnit: SignalUnit) => new SignalUnitState(signalUnit),
	HistorySignalTypeModify
);

export function getSignalUnitState(entityId: string) {
	return provider.get(entityId);
}

export async function loadSignalUnit(entityId: string) {
	const signalUnit = await SignalUnitService.Get(entityId);
	provider.add(signalUnit);
}

export function useSignalUnit(entityId: string) {
	loadSignalUnit(entityId);
}

class SignalUnitState extends EntityState<SignalUnit> {
	// signalUnit = $state<SignalUnit | undefined>();

	constructor(signalUnit: SignalUnit) {
		super(signalUnit);
	}

	reload(entityId: string) {
		loadSignalUnit(entityId);
	}

	updateName(name: string) {
		this.update(SignalUnitService.UpdateName(this.entity.entityId, name));
	}

	updateDesc(desc: string) {
		this.update(SignalUnitService.UpdateDesc(this.entity.entityId, desc));
	}
	
}
