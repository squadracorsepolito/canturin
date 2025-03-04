import { NetworkService, type Network } from '$lib/api/canturin';
import { HistoryNetworkModify } from '$lib/constants/events';
import { EntityState } from '$lib/state/entity-state.svelte';
import { Events as wails } from '@wailsio/runtime';

let networkState: NetworkState;

export function getNetworkState() {
	return networkState;
}

export async function loadNetwork() {
	const network = await NetworkService.Get();
	networkState = new NetworkState(network);
}

class NetworkState extends EntityState<Network> {
	constructor(network: Network) {
		super(network);

		wails.On(HistoryNetworkModify, (e: wails.WailsEvent) => {
			this.entity = e.data[0] as Network;
		});
	}

	updateName(name: string) {
		this.update(NetworkService.UpdateName({ name }));
	}

	updateDesc(desc: string) {
		this.update(NetworkService.UpdateDesc({ desc }));
	}

	addBus() {
		this.update(NetworkService.AddBus());
	}

	deleteBuses(busEntityIDs: string[]) {
		this.update(NetworkService.DeleteBuses({ busEntityIDs }));
	}
}
