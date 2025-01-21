import { SignalService, type Signal } from '$lib/api/canturin';
import { HistorySignalModify } from '$lib/constants/events';
import { EntityState } from '$lib/state/entity-state.svelte';
import { StateProvider } from '$lib/state/state-provider.svelte';

const provider = new StateProvider(
	(signal: Signal) => new SignalState(signal),
	HistorySignalModify
);

export function getSignalState(entityId: string) {
	return provider.get(entityId);
}

export async function loadSignal(entityId: string) {
	const signal = await SignalService.Get(entityId);
	provider.add(signal);
}

export async function deleteSignal(entityId: string) {
	// TODO! implement
	console.log(entityId);
}

class SignalState extends EntityState<Signal> {
	constructor(signal: Signal) {
		super(signal);
	}

	async getInvalidNames() {
		const invalidNames = await SignalService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	updateName(name: string) {
		this.update(SignalService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(SignalService.UpdateDesc(this.entity.entityId, { desc }));
	}

	updateSignalType(signalTypeEntityId: string) {
		this.update(SignalService.UpdateSignalType(this.entity.entityId, { signalTypeEntityId }));
	}

	updateSignalUnit(signalUnitEntityId: string) {
		this.update(SignalService.UpdateSignalUnit(this.entity.entityId, { signalUnitEntityId }));
	}
}
