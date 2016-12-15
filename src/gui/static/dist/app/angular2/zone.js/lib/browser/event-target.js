/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
System.register(['../common/utils'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var utils_1;
    var WTF_ISSUE_555, NO_EVENT_TARGET, EVENT_TARGET;
    function eventTargetPatch(_global) {
        var apis = [];
        var isWtf = _global['wtf'];
        if (isWtf) {
            // Workaround for: https://github.com/google/tracing-framework/issues/555
            apis = WTF_ISSUE_555.split(',').map((v) => 'HTML' + v + 'Element').concat(NO_EVENT_TARGET);
        }
        else if (_global[EVENT_TARGET]) {
            apis.push(EVENT_TARGET);
        }
        else {
            // Note: EventTarget is not available in all browsers,
            // if it's not available, we instead patch the APIs in the IDL that inherit from EventTarget
            apis = NO_EVENT_TARGET;
        }
        for (var i = 0; i < apis.length; i++) {
            var type = _global[apis[i]];
            utils_1.patchEventTargetMethods(type && type.prototype);
        }
    }
    exports_1("eventTargetPatch", eventTargetPatch);
    return {
        setters:[
            function (utils_1_1) {
                utils_1 = utils_1_1;
            }],
        execute: function() {
            WTF_ISSUE_555 = 'Anchor,Area,Audio,BR,Base,BaseFont,Body,Button,Canvas,Content,DList,Directory,Div,Embed,FieldSet,Font,Form,Frame,FrameSet,HR,Head,Heading,Html,IFrame,Image,Input,Keygen,LI,Label,Legend,Link,Map,Marquee,Media,Menu,Meta,Meter,Mod,OList,Object,OptGroup,Option,Output,Paragraph,Pre,Progress,Quote,Script,Select,Source,Span,Style,TableCaption,TableCell,TableCol,Table,TableRow,TableSection,TextArea,Title,Track,UList,Unknown,Video';
            NO_EVENT_TARGET = 'ApplicationCache,EventSource,FileReader,InputMethodContext,MediaController,MessagePort,Node,Performance,SVGElementInstance,SharedWorker,TextTrack,TextTrackCue,TextTrackList,WebKitNamedFlow,Window,Worker,WorkerGlobalScope,XMLHttpRequest,XMLHttpRequestEventTarget,XMLHttpRequestUpload,IDBRequest,IDBOpenDBRequest,IDBDatabase,IDBTransaction,IDBCursor,DBIndex'
                .split(',');
            EVENT_TARGET = 'EventTarget';
        }
    }
});

//# sourceMappingURL=event-target.js.map
