export type SummaryInfo = {
	title: string;
	value: string | number;
	desc?: string;
	badge?: {
		text: string;
		color: 'primary' | 'secondary';
	};
};
