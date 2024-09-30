<script lang="ts">
	import Input from './input.svelte';
	import type { InputProps } from './types';

	let { value = $bindable(), ...restProps }: Omit<InputProps<string>, 'type'> = $props();
</script>

<Input bind:value type="text" {...restProps} />
