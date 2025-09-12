/* SystemJS module definition */
declare var module: NodeModule; // eslint-disable-line no-var
interface NodeModule {
  id: string;
}
/**
 * Needed for using System in app.translate-loader.ts.
 */
declare let System: System; // eslint-disable-line @typescript-eslint/naming-convention
interface System {
  import(request: string): Promise<any>;
}
