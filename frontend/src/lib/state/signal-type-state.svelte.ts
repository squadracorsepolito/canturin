import { SignalTypeService, type SignalType } from '$lib/api/canturin';
import { HistorySignalTypeModify } from '$lib/api/events';
import { EntityState } from './entity-state.svelte';
import { StateProvider } from './state-provider.svelte';

const provider = new StateProvider(
	(signalType: SignalType) => new SignalTypeState(signalType),
	HistorySignalTypeModify
);

export function getSignalTypeState(entityId: string) {
	return provider.get(entityId);
}

export async function loadSignalType(entityId: string) {
	const signalType = await SignalTypeService.Get(entityId);
	provider.add(signalType);
}

class SignalTypeState extends EntityState<SignalType> {
	constructor(signalType: SignalType) {
		super(signalType);
	}

	async getInvalidNames() {
		const invalidNames = await SignalTypeService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	updateName(name: string) {
		this.update(SignalTypeService.UpdateName(this.entity.entityId, name));
	}

	updateDesc(desc: string) {
		this.update(SignalTypeService.UpdateDesc(this.entity.entityId, desc));
	}

	updateMin(min: number) {
		this.update(SignalTypeService.UpdateMin(this.entity.entityId, min));
	}

	updateMax(max: number) {
		this.update(SignalTypeService.UpdateMax(this.entity.entityId, max));
	}

	updateScale(scale: number) {
		this.update(SignalTypeService.UpdateScale(this.entity.entityId, scale));
	}

	updateOffset(offset: number) {
		this.update(SignalTypeService.UpdateOffset(this.entity.entityId, offset));
	}
}
