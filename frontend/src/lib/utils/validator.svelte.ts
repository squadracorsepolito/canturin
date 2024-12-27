import * as v from 'valibot';

export class Validator<T extends v.BaseSchema<unknown, unknown, v.BaseIssue<unknown>>> {
	errors = $state<string[]>();

	constructor(schema: T, getter: () => v.InferInput<typeof schema>) {
		$effect(() => {
			const res = v.safeParse(schema, getter());
			if (res.success) {
				this.errors = undefined;
				return;
			}

			this.errors = v.flatten(res.issues).root;
		});
	}
}

export function checkUnused<T>(getter: () => T[]) {
	return v.check(
		(item: T) => !getter().includes(item),
		(issue) => `Invalid value: Received ${issue.received} but already used`
	);
}
