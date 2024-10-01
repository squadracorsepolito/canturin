<script lang="ts">
	import Input from './input.svelte';
	import type { TextInputProps } from './types';

	let { value = $bindable(), ...restProps }: TextInputProps = $props();
</script>

<Input bind:value type="text" {...restProps} />
