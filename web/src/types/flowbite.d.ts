declare module 'flowbite-datepicker' {
  export class Datepicker {
    constructor(targetElement: HTMLElement, options?: DatePickerOption);

    element: HTMLElement;
    dates: Array<number>;

    hide(): void;
    destroy(): void;
  }

  export type DatePickerOption = {
    format?: string;
    clearBtn?: boolean;
    autohide?: boolean;
  };
}

interface IDatePickerChangeEvent {
  date: Date;
  datepicker: import('flowbite-datepicker').Datepicker;
  viewDate: Date;
  viewId: number;
}
