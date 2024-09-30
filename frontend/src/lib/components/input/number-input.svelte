<script lang="ts">
	import type { InputProps } from './types';
	import Input from './input.svelte';

	let { value = $bindable(), ...restProps }: Omit<InputProps<number>, 'type'> = $props();
</script>

<Input bind:value type="number" {...restProps} />
