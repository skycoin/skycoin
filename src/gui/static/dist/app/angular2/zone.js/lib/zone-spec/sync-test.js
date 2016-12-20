/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
class SyncTestZoneSpec {
    constructor(namePrefix) {
        this.runZone = Zone.current;
        this.name = 'syncTestZone for ' + namePrefix;
    }
    onScheduleTask(delegate, current, target, task) {
        switch (task.type) {
            case 'microTask':
            case 'macroTask':
                throw new Error(`Cannot call ${task.source} from within a sync test.`);
            case 'eventTask':
                task = delegate.scheduleTask(target, task);
                break;
        }
        return task;
    }
}
// Export the class so that new instances can be created with proper
// constructor params.
Zone['SyncTestZoneSpec'] = SyncTestZoneSpec;

//# sourceMappingURL=sync-test.js.map
