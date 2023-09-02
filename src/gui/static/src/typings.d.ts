/* eslint-disable no-var */

/* SystemJS module definition */
declare var module: NodeModule;
interface NodeModule {
  id: string;
}

declare var System: System;
interface System {
  import(request: string): Promise<any>;
}
