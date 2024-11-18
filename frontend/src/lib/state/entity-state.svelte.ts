import { pushToast } from '$lib/components/toast/toast-provider.svelte';

export type Entity = {
	entityId: string;
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
			this.#fallback = this.entity;
			this.entity = newEntity;
		} catch (error) {
			this.entity = this.#fallback;
			pushToast('error', 'Error', 'Operation failed');
			console.error(error);
		}
	}
}
