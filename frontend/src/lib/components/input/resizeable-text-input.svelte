<script lang="ts">
	import Input from './input.svelte';
	import type { ResizeableTextInputProps } from './types';

	let { value = $bindable(), size = 'md', ...restProps }: ResizeableTextInputProps = $props();

	let width = $state(0);

	function getDummyClass() {
		let s = '';
		switch (size) {
			case 'md':
				s = 'px-4 text-sm';
				break;
			case 'sm':
				s = 'px-3 text-sm';
				break;
		}

		return `absolute break-all invisible ${s}`;
	}
</script>

<div class="relative">
	<Input bind:value type="text" {size} width={width + 6} {...restProps} />

	<div bind:clientWidth={width} class={getDummyClass()}>
		{value}
	</div>
</div>
