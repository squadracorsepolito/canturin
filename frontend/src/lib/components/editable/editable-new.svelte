<script lang="ts" generics="T extends AnyZodObject">
	import { FiniteStateMachine } from 'runed';
	import type { Snippet } from 'svelte';
	import type { AnyZodObject, z } from 'zod';

	type Values = z.infer<T>;
	type Errors = z.inferFlattenedErrors<T>['fieldErrors'];
	type InputSnippetProps = {
		fsm: typeof fsm;
		values: Values;
		errors: Errors;
	};

	type Props = {
		hidePlaceholder?: boolean;
		placeholder: Snippet;
		input: Snippet<[InputSnippetProps]>;
		schema: T;
		initialValues: Values;
		onsubmit: (values: Values) => void;
	};

	let { hidePlaceholder, placeholder, input, initialValues, schema, onsubmit }: Props = $props();

	type State = 'idle' | 'editing' | 'resetting';
	type Events = 'DBLCLICK' | 'ESCAPE' | 'BLUR' | 'SUBMIT' | 'TIMEOUT';
	const fsm = new FiniteStateMachine<State, Events>('idle', {
		idle: {
			DBLCLICK: 'editing'
		},
		editing: {
			BLUR: () => {
				if (hasErrors()) {
					return 'resetting';
				}

				onsubmit(values);
				return 'idle';
			},
			// SUBMIT: () => {
			// 	if (!hasErrors()) {
			// 		return 'idle';
			// 	}
			// },
			ESCAPE: 'resetting'
		},
		resetting: {
			_enter: () => {
				fsm.debounce(100, 'TIMEOUT');
			},
			TIMEOUT: 'idle'
		}
	});

	let values = $state(initialValues);
	let errors = $derived.by<Errors>(() => {
		const res = schema.safeParse(values);
		if (res.success) {
			return {};
		}
		return res.error.flatten().fieldErrors;
	});

	function hasErrors() {
		return Object.keys(errors).length > 0;
	}

	function handleSubmit(e: SubmitEvent) {
		if (hasErrors()) {
			e.preventDefault();
			return;
		}

		// fsm.send('SUBMIT');
		onsubmit(values);
	}
</script>

{#if (hidePlaceholder && fsm.current === 'idle') || !hidePlaceholder}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div ondblclick={() => fsm.send('DBLCLICK')}>
		{@render placeholder()}
	</div>
{/if}

{#if fsm.current === 'editing'}
	<form onsubmit={handleSubmit}>
		{@render input({ fsm, values, errors })}
	</form>
{/if}
