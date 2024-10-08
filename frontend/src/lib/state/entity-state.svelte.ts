export type Entity = {
	entityId: string;
};

export class EntityState<E extends Entity> {
	entity = $state() as E;

	constructor(entity: E) {
		this.entity = entity;
	}

	async update(promise: Promise<E>) {
		try {
			this.entity = await promise;
		} catch (error) {
			console.error(error);
		}
	}
}
