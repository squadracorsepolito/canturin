<script lang="ts" generics="T extends number|string">
	import type { Action } from 'svelte/action';
	import type { InputProps } from './types';

	let {
		value = $bindable(),
		name,
		type,
		label,
		errors,
		focusOnDiplay,
		onblur,
		onescape
	}: InputProps<T> = $props();

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

<label for={name} class="form-control max-w-xs">
	{#if label}
		<div class="label">
			<span class="label-text-alt">
				{label}
			</span>
		</div>
	{/if}

	<input
		bind:value
		{type}
		{name}
		onblur={handleBlur}
		onkeydown={handleKeydown}
		use:action
		class="input {errors ? 'input-error bg-error text-error-content' : 'input-primary'}"
	/>

	<div class="label h-8">
		{#if errors}
			{#each errors as err}
				<span class="label-text-alt text-error">{err}</span>
			{/each}
		{/if}
	</div>
</label>
