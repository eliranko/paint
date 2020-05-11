export class Canvas {
    uuid: string = "";
    name: string;
    data: string;

    constructor(name: string, data: string) {
        this.name = name;
        this.data = data;
    }
}