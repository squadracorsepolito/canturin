import plugin from 'tailwindcss/plugin';

const colors = {
	primary: '--p',
	'primary-content': '--pc',

	secondary: '--s',
	'secondary-content': '--sc',

	accent: '--a',
	'accent-content': '--ac',

	neutral: '--n',
	'neutral-content': '--nc',

	'base-100': '--b1',
	'base-200': '--b2',
	'base-300': '--b3',
	'base-content': '--bc',

	info: '--in',
	'info-content': '--inc',

	success: '--su',
	'success-content': '--suc',

	warning: '--wa',
	'warning-content': '--wac',

	error: '--er',
	'error-content': '--erc'
};

export default plugin(({ addUtilities }) => {
	for (const [color, cssVar] of Object.entries(colors)) {
		const className = `.bg-${color}-ghost`;

		addUtilities({
			[className]: {
				backgroundColor: `var(--fallback-p,oklch(var(${cssVar})/0.05))`
			}
		});
	}

	for (const [color, cssVar] of Object.entries(colors)) {
		const className = `.focus-ring-${color}`;

		addUtilities({
			[className]: {
				'--tw-ring-offset-shadow':
					'var(--tw-ring-inset) 0 0 0 var(--tw-ring-offset-width) var(--tw-ring-offset-color)',
				'--tw-ring-shadow':
					'var(--tw-ring-inset) 0 0 0 calc(2px + var(--tw-ring-offset-width)) var(--tw-ring-color)',
				'box-shadow':
					'var(--tw-ring-offset-shadow), var(--tw-ring-shadow), var(--tw-shadow, 0 0 #0000)',
				'--tw-ring-color': `var(--fallback-p,oklch(var(${cssVar})/0.25))`
			}
		});
	}
});
