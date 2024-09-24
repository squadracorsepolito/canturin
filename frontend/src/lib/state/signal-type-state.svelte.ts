import { SignalTypeService, type SignalType } from '$lib/api/canturin';

import { EntityState } from './entity-state.svelte';

// async function loadSignalType(state: SignalTypeState, entityId: string) {
// 	state.isLoading = true;

// 	try {
// 		const sigType = await SignalTypeService.GetOpen(entityId);
// 		state.signalType = sigType;
// 	} catch (error) {
// 		state.signalType = undefined;
// 		console.error(error);
// 	}

// 	state.isLoading = false;
// }

// class SignalTypeState extends AsyncState {
// 	signalType = $state<SignalType>();

// 	reload(entityId: string) {
// 		loadSignalType(this, entityId);
// 	}

// 	async update(promise: Promise<SignalType>) {
// 		try {
// 			const sigType = await promise;
// 			this.signalType = sigType;
// 		} catch (error) {
// 			console.error(error);
// 		}
// 	}

// 	async updateName(newName: string) {
// 		if (!this.signalType) return;

// 		await this.update(SignalTypeService.UpdateName(this.signalType.entityId, newName));
// 	}

// 	async updateDesc(newDesc: string) {
// 		if (!this.signalType) return;

// 		await this.update(SignalTypeService.UpdateDesc(this.signalType.entityId, newDesc));
// 	}
// }

// export function useSignalType(entityId: string) {
// 	const state = new SignalTypeState();

// 	loadSignalType(state, entityId);

// 	return state;
// }

class SignalTypeState extends EntityState<SignalType> {
	constructor() {
		super(SignalTypeService.GetOpen);
	}

	updateName(newName: string) {
		if (!this.entity) return;

		this.update(SignalTypeService.UpdateName(this.entity.entityId, newName));
	}

	updateDesc(newDesc: string) {
		if (!this.entity) return;

		this.update(SignalTypeService.UpdateDesc(this.entity.entityId, newDesc));
	}
}

export function useSignalType(entityId: string) {
	const state = new SignalTypeState();
	state.load(entityId);
	return state;
}
