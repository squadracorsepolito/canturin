import { NetworkService, type Network } from '$lib/api/canturin';
import { HistoryNetworkModify, NetworkLoaded } from '$lib/constants/events';
import { EntityState } from '$lib/state/entity-state.svelte';
import { Events as wails } from '@wailsio/runtime';

let state: NetworkState;

let isLoaded = $state(false);

wails.On(NetworkLoaded, () => {
	getNetwork();
});

export function getNetworkState() {
	return state;
}

export function isNetworkLoaded() {
	return isLoaded;
}

export async function getNetwork() {
	const network = await NetworkService.Get();
	state = new NetworkState(network);
}

export async function loadNetwork(path: string) {
	await NetworkService.Load(path);
}

export async function createNetwork() {
	await NetworkService.Create();
}

class NetworkState extends EntityState<Network> {
	constructor(network: Network) {
		super(network);

		isLoaded = true;

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
