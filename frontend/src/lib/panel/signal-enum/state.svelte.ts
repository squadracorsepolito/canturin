import { SignalEnumService, type SignalEnum } from '$lib/api/canturin';
import { HistorySignalEnumModify } from '$lib/constants/events';
import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import layout from '$lib/state/layout-state.svelte';
import { EntityState } from '../../state/entity-state.svelte';
import { StateProvider } from '../../state/state-provider.svelte';

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

export async function deleteSignalEnum(entityId: string) {
	try {
		await SignalEnumService.Delete(entityId);
		provider.remove(entityId);
		layout.closeIfOpen(entityId);
	} catch (error) {
		pushToast('error', 'Error', 'Operation failed');
		console.error(error);
	}
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
		this.update(SignalEnumService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(SignalEnumService.UpdateDesc(this.entity.entityId, { desc }));
	}

	reorderValue(valueEntityId: string, from: number, to: number) {
		if (from === to) return;

		this.update(
			SignalEnumService.ReorderValue(this.entity.entityId, {
				valueEntityId: valueEntityId,
				from,
				to
			})
		);
	}

	addValue() {
		this.update(SignalEnumService.AddValue(this.entity.entityId));
	}

	deleteValue(valueEntityId: string) {
		this.update(
			SignalEnumService.RemoveValues(this.entity.entityId, {
				valueEntityIds: [valueEntityId]
			})
		);
	}

	deleteValues(valueEntityIds: string[]) {
		this.update(
			SignalEnumService.RemoveValues(this.entity.entityId, {
				valueEntityIds
			})
		);
	}

	updateValueName(valueEntityId: string, name: string) {
		this.update(
			SignalEnumService.UpdateValueName(this.entity.entityId, {
				valueEntityId,
				name
			})
		);
	}

	updateValueDesc(valueEntityId: string, desc: string) {
		this.update(
			SignalEnumService.UpdateValueDesc(this.entity.entityId, {
				valueEntityId,
				desc
			})
		);
	}

	updateValueIndex(valueEntityId: string, index: number) {
		this.update(
			SignalEnumService.UpdateValueIndex(this.entity.entityId, {
				valueEntityId,
				index
			})
		);
	}
}
