<script lang="ts">
	import type { Action } from 'svelte/action';
	import { ZodError, ZodNumber } from 'zod';

	type Props = {
		value: number;
		name: string;
		label?: string;
		errors?: string[];
		focusOnDiplay?: boolean;
		onblur?: () => void;
		onescape?: () => void;
	};

	let {
		value = $bindable(),
		name,
		label,
		errors,
		focusOnDiplay,
		onblur,
		onescape
	}: Props = $props();

	function handleBlur() {
		if (!onblur) return;

		onblur();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!onescape) return;

		if (e.key === 'Escape') onescape();
	}

	const action: Action<HTMLInputElement> = (el) => {
		if (focusOnDiplay) {
			el.focus();
		}
	};
</script>

<div>
	<label for={name}>
		{#if label}
			{label}
		{/if}
	</label>

	<input
		bind:value
		type="number"
		{name}
		onblur={handleBlur}
		onkeydown={handleKeydown}
		use:action
		class="input input-primary"
	/>

	{#if errors}
		{#each errors as err}
			<span>{err}</span>
		{/each}
	{/if}
</div>
