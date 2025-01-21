import { BusService, BusType, type Bus } from '$lib/api/canturin';
import { HistoryBusModify } from '$lib/constants/events';
import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import { EntityState } from '$lib/state/entity-state.svelte';
import layout from '$lib/state/layout-state.svelte';
import { StateProvider } from '$lib/state/state-provider.svelte';

const provider = new StateProvider((bus: Bus) => new BusState(bus), HistoryBusModify);

export function getBusState(entityId: string) {
	return provider.get(entityId);
}

export async function loadBus(entityId: string) {
	const bus = await BusService.Get(entityId);
	provider.add(bus);
}

export async function deleteBus(entityId: string) {
	try {
		await BusService.Delete(entityId);
		provider.remove(entityId);
		layout.closeIfOpen(entityId);
	} catch (error) {
		pushToast('error', 'Error', 'Operation failed');
		console.error(error);
	}
}

class BusState extends EntityState<Bus> {
	constructor(bus: Bus) {
		super(bus);
	}

	async getInvalidNames() {
		const invalidNames = await BusService.GetInvalidNames(this.entity.entityId);

		if (invalidNames) {
			return invalidNames;
		}

		return [];
	}

	updateName(name: string) {
		this.update(BusService.UpdateName(this.entity.entityId, { name }));
	}

	updateDesc(desc: string) {
		this.update(BusService.UpdateDesc(this.entity.entityId, { desc }));
	}

	updateBusType(busType: BusType) {
		this.update(BusService.UpdateBusType(this.entity.entityId, { busType }));
	}

	updateBaudrate(baudrate: number) {
		this.update(BusService.UpdateBaudrate(this.entity.entityId, { baudrate }));
	}
}
