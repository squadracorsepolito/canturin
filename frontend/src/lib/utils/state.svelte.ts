export class DerivedBindable<T> {
	#value = $state() as T;

	constructor(getter: () => T) {
		$effect(() => {
			this.#value = getter();
		});
	}

	get value() {
		return this.#value;
	}

	set value(value: T) {
		this.#value = value;
	}
}
