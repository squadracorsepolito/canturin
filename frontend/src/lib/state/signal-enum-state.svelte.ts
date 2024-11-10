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
	indexes = $state<number[]>([]);

	constructor(signalEnum: SignalEnum) {
		super(signalEnum);

		this.getIndexes(signalEnum);
	}
	private getIndexes(sigEnum: SignalEnum) {
		if (!sigEnum.values) return;

		this.indexes = sigEnum.values.map((val) => val.index);
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

	reorderValue(valueEntID: string, from: number, to: number) {
		this.update(SignalEnumService.ReorderValue(this.entity.entityId, valueEntID, from, to));
	}

	addValue() {
		this.update(SignalEnumService.AddValue(this.entity.entityId));
	}

	deleteValue(valueEntID: string) {
		this.update(SignalEnumService.RemoveValues(this.entity.entityId, valueEntID));
	}

	deleteValues(valueEntIds: string[]) {
		this.update(SignalEnumService.RemoveValues(this.entity.entityId, ...valueEntIds));
	}

	updateValueName(valueEntID: string, name: string) {
		this.update(SignalEnumService.UpdateValueName(this.entity.entityId, valueEntID, name));
	}

	updateValueIndex(valueEntID: string, index: number) {
		const f = async () => {
			const sigEnum = await SignalEnumService.UpdateValueIndex(
				this.entity.entityId,
				valueEntID,
				index
			);
			this.getIndexes(sigEnum);
			return sigEnum;
		};

		this.update(f());
	}

	updateValueDesc(valueEntID: string, desc: string) {
		this.update(SignalEnumService.UpdateValueDesc(this.entity.entityId, valueEntID, desc));
	}
}
