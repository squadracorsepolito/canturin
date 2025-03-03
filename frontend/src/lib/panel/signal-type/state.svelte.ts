import { SignalTypeKind, SignalTypeService, type SignalType } from '$lib/api/canturin';
import { HistorySignalTypeModify } from '$lib/constants/events';
import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import { EntityState } from '../../state/entity-state.svelte';
import layout from '../../state/layout-state.svelte';
import { StateProvider } from '../../state/state-provider.svelte';

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

export async function createSignalType(signalTypeKind: SignalTypeKind, size: number) {
	try {
		await SignalTypeService.Create({ signalTypeKind, size });
	} catch (err) {
		console.error(err);
		pushToast('error', 'Error', 'Operation failed');
	}
}

export async function deleteSignalType(entityId: string) {
	try {
		await SignalTypeService.Delete(entityId);
		provider.remove(entityId);
		layout.closeIfOpen(entityId);
	} catch (err) {
		console.error(err);
		pushToast('error', 'Error', 'Operation failed');
	}
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
		this.update(SignalTypeService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(SignalTypeService.UpdateDesc(this.entity.entityId, { desc }));
	}

	updateSigned(signed: boolean) {
		this.update(SignalTypeService.UpdateSigned(this.entity.entityId, { signed }));
	}

	updateMin(min: number) {
		this.update(SignalTypeService.UpdateMin(this.entity.entityId, { min }));
	}

	updateMax(max: number) {
		this.update(SignalTypeService.UpdateMax(this.entity.entityId, { max }));
	}

	updateScale(scale: number) {
		this.update(SignalTypeService.UpdateScale(this.entity.entityId, { scale }));
	}

	updateOffset(offset: number) {
		this.update(SignalTypeService.UpdateOffset(this.entity.entityId, { offset }));
	}
}
