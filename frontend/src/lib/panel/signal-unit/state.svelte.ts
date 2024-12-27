import { SignalUnitKind, SignalUnitService, type SignalUnit } from '$lib/api/canturin';
import { HistorySignalUnitModify } from '$lib/api/events';
import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import layout from '$lib/state/layout-state.svelte';
import { EntityState } from '../../state/entity-state.svelte';
import { StateProvider } from '../../state/state-provider.svelte';

const provider = new StateProvider(
	(signalUnit: SignalUnit) => new SignalUnitState(signalUnit),
	HistorySignalUnitModify
);

export function getSignalUnitState(entityId: string) {
	return provider.get(entityId);
}

export async function loadSignalUnit(entityId: string) {
	const signalUnit = await SignalUnitService.Get(entityId);
	provider.add(signalUnit);
}

export async function deleteSignalUnit(entityId: string) {
	try {
		await SignalUnitService.Delete(entityId);
		provider.remove(entityId);
		layout.closeIfOpen(entityId);
	} catch (error) {
		pushToast('error', 'Error', 'Operation failed');
		console.error(error);
	}
}

class SignalUnitState extends EntityState<SignalUnit> {
	constructor(signalUnit: SignalUnit) {
		super(signalUnit);
	}

	async getInvalidNames() {
		const invalidNames = await SignalUnitService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	updateName(name: string) {
		this.update(SignalUnitService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(SignalUnitService.UpdateDesc(this.entity.entityId, { desc }));
	}

	updateKind(kindStr: string) {
		const kind = kindStr as SignalUnitKind;
		this.update(SignalUnitService.UpdateKind(this.entity.entityId, { kind }));
	}

	updateSymbol(symbol: string) {
		this.update(SignalUnitService.UpdateSymbol(this.entity.entityId, { symbol }));
	}
}
