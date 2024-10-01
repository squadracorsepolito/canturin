<script lang="ts">
	import type { NumberInputProps } from './types';
	import Input from './input.svelte';

	let { value = $bindable(), ...restProps }: NumberInputProps = $props();
</script>

<Input bind:value type="number" {...restProps} />
