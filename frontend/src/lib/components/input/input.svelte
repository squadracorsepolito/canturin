<script lang="ts" generics="T extends number|string">
	import type { Action } from 'svelte/action';
	import type { InputProps } from './types';

	let {
		value = $bindable(),
		name,
		type,
		label,
		errors,
		focusOnDisplay,
		size = 'md',
		width,
		min,
		max,
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
		if (focusOnDisplay) {
			el.focus();
		}
	};

	function getInputClass() {
		let s = '';
		switch (size) {
			case 'md':
				s = 'input-md';
				break;
			case 'sm':
				s = 'input-sm';
				break;
		}

		let color = 'input-primary';
		if (errors) {
			color = 'input-error bg-error text-error-content';
		}

		return `input ${s} ${color}`;
	}
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
		{min}
		{max}
		onblur={handleBlur}
		onkeydown={handleKeydown}
		use:action
		class={getInputClass()}
		style:width={width && `${width}px`}
		step="any"
	/>
</label>

{#if errors}
	<div class="pt-1 absolute">
		{#each errors as err}
			<span class="label-text-alt text-error">{err}</span>
		{/each}
	</div>
{/if}
