import { SignalEnumService, type SignalEnum } from '$lib/api/canturin';
import { HistorySignalEnumModify } from '$lib/api/events';
import { EntityState } from './entity-state.svelte';
import { StateProvider } from './state-provider.svelte';

const provider = new StateProvider(
	(signalEnum: SignalEnum) => new SignalEnumState(signalEnum),
	HistorySignalEnumModify
);

export function getSignalEnumState(entityId: string) {
	return provider.get(entityId);
}

export async function loadSignalEnum(entityId: string) {
	const signalEnum = await SignalEnumService.Get(entityId);
	provider.add(signalEnum);
}

class SignalEnumState extends EntityState<SignalEnum> {
	constructor(signalEnum: SignalEnum) {
		super(signalEnum);
	}

	async getInvalidNames() {
		const invalidNames = await SignalEnumService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	updateName(name: string) {
		this.update(SignalEnumService.UpdateName(this.entity.entityId, name));
	}

	updateDesc(desc: string) {
		this.update(SignalEnumService.UpdateDesc(this.entity.entityId, desc));
	}

	updateValueName(valueEntID: string, name: string) {
		this.update(SignalEnumService.UpdateValueName(this.entity.entityId, valueEntID, name));
	}

	updateValueIndex(valueEntID: string, index: number) {
		this.update(SignalEnumService.UpdateValueIndex(this.entity.entityId, valueEntID, index));
	}

	updateValueDesc(valueEntID: string, desc: string) {
		this.update(SignalEnumService.UpdateValueDesc(this.entity.entityId, valueEntID, desc));
	}
}
