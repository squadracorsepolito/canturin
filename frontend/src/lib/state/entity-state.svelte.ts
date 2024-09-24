export class EntityState<
	T extends {
		entityId: string;
	}
> {
	loadFn: (entityId: string) => Promise<T>;

	isLoading = $state(false);
	entity = $state<T>();

	constructor(loadFn: (entityId: string) => Promise<T>) {
		this.loadFn = loadFn;
	}

	async load(entityId: string) {
		this.isLoading = true;

		try {
			this.entity = await this.loadFn(entityId);
		} catch (error) {
			this.entity = undefined;
			console.error(error);
		} finally {
			this.isLoading = false;
		}
	}

	reload(entityId: string) {
		this.load(entityId);
	}

	async update(updatePromise: Promise<T>) {
		try {
			this.entity = await updatePromise;
		} catch (error) {
			console.error(error);
		}
	}
}
