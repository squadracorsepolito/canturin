<script lang="ts">
	import { selectTextOnFocus } from '$lib/actions';
	import type { Action } from 'svelte/action';
	import { z, ZodError } from 'zod';

	type Props = {
		prefixName: string;
		initialValue: string;
		invalidNames: string[];
		onsubmitname: (name: string) => void;
	};

	let { prefixName, initialValue, invalidNames, onsubmitname }: Props = $props();

	const schema = z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated name' })
	});

	let value = $state(initialValue);
	let width = $state(0);
	let error = $derived.by(() => {
		try {
			schema.parse({ name: value });
		} catch (err) {
			if (err instanceof ZodError) {
				return err.issues[0].message;
			}
		}
	});

	function resetValue() {
		value = initialValue;
	}

	function handleSubmit() {
		if (value !== initialValue) {
			onsubmitname(value);
		}
	}

	function onsubmit(e: SubmitEvent) {
		if (error) {
			e.preventDefault();
			return;
		}

		handleSubmit();
	}

	const inputAction: Action<HTMLInputElement> = (node) => {
		const handleBlur = () => {
			if (!node) return;

			if (error) {
				resetValue();
				return;
			}

			handleSubmit();
		};

		const handleKeydown = (e: KeyboardEvent) => {
			if (!node) return;

			if (e.key === 'Escape') {
				resetValue();
				node.blur();
			}
		};

		const handleFocus = () => {
			if (node && typeof node.select === 'function') {
				node.select();
			}
		};

		node.addEventListener('blur', handleBlur);
		node.addEventListener('keydown', handleKeydown);
		node.addEventListener('focus', handleFocus);

		return {
			destroy() {
				node.removeEventListener('blur', handleBlur);
				node.removeEventListener('keydown', handleKeydown);
				node.removeEventListener('focus', handleFocus);
			}
		};
	};
</script>

<form class="relative" {onsubmit}>
	<input
		type="text"
		name={prefixName + '_name'}
		bind:value
		class="input {error
			? 'input-error'
			: 'input-primary'} font-medium text-xl border-transparent focus:border-solid"
		style:width={`${width + 2 + 1}px`}
		use:inputAction
	/>

	{#if error}
		<span class="text-sm text-error">{error}</span>
	{/if}

	<div
		bind:clientWidth={width}
		class="px-4 inline-block absolute left-0 break-all top-12 font-medium text-xl invisible"
	>
		{value}
	</div>
</form>
