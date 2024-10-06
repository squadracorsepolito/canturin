<script lang="ts" generics="T extends AnyZodObject">
	import { clickOutside } from '$lib/actions';
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
		blurOnSubmit?: boolean;
		submitOnClickOutside?: boolean;
		schema: T;
		initialValues: Values;
		onsubmit: (values: Values) => void;
		placeholder: Snippet<[typeof fsm]>;
		input: Snippet<[InputSnippetProps]>;
	};

	let {
		hidePlaceholder,
		blurOnSubmit,
		submitOnClickOutside,
		placeholder,
		input,
		initialValues,
		schema,
		onsubmit
	}: Props = $props();

	type State = 'idle' | 'editing' | 'resetting';
	type Events = 'DBLCLICK' | 'ESCAPE' | 'BLUR' | 'TIMEOUT';
	const fsm = new FiniteStateMachine<State, Events>('idle', {
		idle: {
			DBLCLICK: 'editing',
			BLUR: () => {}
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
			},
			BLUR: () => {}
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

		if (blurOnSubmit) {
			fsm.send('BLUR');
		}
	}

	function handleClickOutside() {
		if (submitOnClickOutside) {
			fsm.send('BLUR');
		}
	}
</script>

<div use:clickOutside={handleClickOutside}>
	{#if (hidePlaceholder && fsm.current === 'idle') || !hidePlaceholder}
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div ondblclick={() => fsm.send('DBLCLICK')}>
			{@render placeholder(fsm)}
		</div>
	{/if}

	{#if fsm.current !== 'idle'}
		<form onsubmit={handleSubmit}>
			{@render input({ fsm, values, errors })}
		</form>
	{/if}
</div>
