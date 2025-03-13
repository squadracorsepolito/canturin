import { pushToast } from '$lib/components/toast/toast-provider.svelte';
import { updatePanelName } from '$lib/panel/panel-stack-state.svelte';

export type Entity = {
	entityId: string;
	name: string;
};

export class EntityState<E extends Entity> {
	#fallback: E;
	entity = $state() as E;

	constructor(entity: E) {
		this.#fallback = entity;
		this.entity = entity;
	}

	async update(promise: Promise<E>) {
		try {
			const newEntity = await promise;
			this.set(newEntity);
		} catch (error) {
			this.entity = this.#fallback;
			pushToast('error', 'Error', 'Operation failed');
			console.error(error);
		}
	}

	set(entity: E) {
		if (entity.name !== this.#fallback.name) {
			updatePanelName(entity.entityId, entity.name);
		}

		this.#fallback = this.entity;
		this.entity = entity;
	}
}
