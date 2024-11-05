<script lang="ts">
	import type { IconProps } from './types';

	let { height = 24, width = 24 }: IconProps = $props();
</script>

<svg xmlns="http://www.w3.org/2000/svg" {width} {height} viewBox="0 0 24 24">
	<g
		fill="none"
		stroke="currentColor"
		stroke-linecap="round"
		stroke-linejoin="round"
		stroke-width="1.5"
	>
		<path d="m4 12.9l3.143 3.6L15 7.5" opacity="0.5" />
		<path d="m20 7.563l-8.571 9L11 16" />
	</g>
</svg>
