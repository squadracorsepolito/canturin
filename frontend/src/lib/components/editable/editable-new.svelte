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
		schema: T;
		initialValues: Values;
		onsubmit: (values: Values) => void;
		placeholder: Snippet;
		input: Snippet<[InputSnippetProps]>;
	};

	let { hidePlaceholder, placeholder, input, initialValues, schema, onsubmit }: Props = $props();

	type State = 'idle' | 'editing' | 'resetting';
	type Events = 'DBLCLICK' | 'ESCAPE' | 'BLUR' | 'TIMEOUT';
	const fsm = new FiniteStateMachine<State, Events>('idle', {
		idle: {
			DBLCLICK: 'editing'
		},
		editing: {
			BLUR: () => {
				if (hasErrors()) {
					return 'resetting';
				}
				if (!isSubmitting) {
					onsubmit(values);
				}
				return 'idle';
			},
			ESCAPE: 'resetting',
			_exit: () => {
				isSubmitting = false;
			}
		},
		resetting: {
			_enter: () => {
				fsm.debounce(100, 'TIMEOUT');
			},
			TIMEOUT: () => {
				values = initialValues;
				return 'idle';
			}
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
	let isSubmitting = $state(false);

	function hasErrors() {
		return Object.keys(errors).length > 0;
	}

	function handleSubmit(e: SubmitEvent) {
		if (hasErrors()) {
			e.preventDefault();
			return;
		}

		isSubmitting = true;
		onsubmit(values);
	}
</script>

{#if (hidePlaceholder && fsm.current === 'idle') || !hidePlaceholder}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div ondblclick={() => fsm.send('DBLCLICK')}>
		{@render placeholder()}
	</div>
{/if}

{#if fsm.current !== 'idle'}
	<form onsubmit={handleSubmit}>
		{@render input({ fsm, values, errors })}
	</form>
{/if}
