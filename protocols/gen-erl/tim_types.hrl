-ifndef(_tim_types_included).
-define(_tim_types_included, yeah).

%% struct 'TimError'

-record('TimError', {'errCode' :: integer(),
                     'errMsg' :: string() | binary()}).
-type 'TimError'() :: #'TimError'{}.

%% struct 'TimNode'

-record('TimNode', {'key' :: string() | binary(),
                    'value' :: string() | binary()}).
-type 'TimNode'() :: #'TimNode'{}.

%% struct 'TimAckBean'

-record('TimAckBean', {'id' :: string() | binary(),
                       'ackType' :: string() | binary(),
                       'ackStatus' :: string() | binary(),
                       'extraList' :: list(),
                       'err' :: 'TimError'(),
                       'extraMap' :: dict:dict()}).
-type 'TimAckBean'() :: #'TimAckBean'{}.

%% struct 'TimHBean'

-record('TimHBean', {'chl' :: integer(),
                     'platform' :: integer(),
                     'version' :: integer()}).
-type 'TimHBean'() :: #'TimHBean'{}.

%% struct 'TimParam'

-record('TimParam', {'timestamp' :: string() | binary(),
                     'version' :: integer(),
                     'lang' :: string() | binary(),
                     'extraList' :: list(),
                     'extraMap' :: dict:dict(),
                     'interflow' :: string() | binary(),
                     'tls' :: string() | binary()}).
-type 'TimParam'() :: #'TimParam'{}.

%% struct 'TimTime'

-record('TimTime', {'timestamp' :: string() | binary(),
                    'formattime' :: string() | binary()}).
-type 'TimTime'() :: #'TimTime'{}.

%% struct 'TimArea'

-record('TimArea', {'country' :: string() | binary(),
                    'province' :: string() | binary(),
                    'city' :: string() | binary(),
                    'extraList' :: list(),
                    'extraMap' :: dict:dict()}).
-type 'TimArea'() :: #'TimArea'{}.

%% struct 'Tid'

-record('Tid', {'name' :: string() | binary(),
                'domain' :: string() | binary(),
                'resource' :: string() | binary(),
                'type' :: string() | binary(),
                'extraList' :: list(),
                'extraMap' :: dict:dict()}).
-type 'Tid'() :: #'Tid'{}.

%% struct 'TimUserBean'

-record('TimUserBean', {'tid' :: 'Tid'(),
                        'nickname' :: string() | binary(),
                        'remarkname' :: string() | binary(),
                        'brithday' :: string() | binary(),
                        'gender' :: integer(),
                        'headurl' :: string() | binary(),
                        'area' :: 'TimArea'(),
                        'headbyte' :: string() | binary(),
                        'photoBytes' :: list(),
                        'extraList' :: list(),
                        'extraMap' :: dict:dict()}).
-type 'TimUserBean'() :: #'TimUserBean'{}.

%% struct 'TimRoom'

-record('TimRoom', {'tid' :: 'Tid'(),
                    'founderTid' :: 'Tid'(),
                    'HostsTid' :: list(),
                    'membersTid' :: list(),
                    'headurl' :: string() | binary(),
                    'roomName' :: string() | binary(),
                    'desc' :: string() | binary(),
                    'createTime' :: 'TimTime'(),
                    'extraList' :: list(),
                    'extraMap' :: dict:dict()}).
-type 'TimRoom'() :: #'TimRoom'{}.

%% struct 'TimPBean'

-record('TimPBean', {'threadId' :: string() | binary(),
                     'fromTid' :: 'Tid'(),
                     'toTid' :: 'Tid'(),
                     'status' :: string() | binary(),
                     'type' :: string() | binary(),
                     'priority' :: integer(),
                     'show' :: string() | binary(),
                     'leaguerTid' :: 'Tid'(),
                     'extraList' :: list(),
                     'error' :: 'TimError'(),
                     'extraMap' :: dict:dict()}).
-type 'TimPBean'() :: #'TimPBean'{}.

%% struct 'TimMBean'

-record('TimMBean', {'threadId' :: string() | binary(),
                     'mid' :: string() | binary(),
                     'fromTid' :: 'Tid'(),
                     'toTid' :: 'Tid'(),
                     'body' :: string() | binary(),
                     'type' :: string() | binary(),
                     'msgType' :: integer(),
                     'offline' :: 'TimTime'(),
                     'leaguerTid' :: 'Tid'(),
                     'extraList' :: list(),
                     'timestamp' :: string() | binary(),
                     'error' :: 'TimError'(),
                     'extraMap' :: dict:dict(),
                     'readstatus' :: integer()}).
-type 'TimMBean'() :: #'TimMBean'{}.

%% struct 'TimIqBean'

-record('TimIqBean', {'threadId' :: string() | binary(),
                      'fromTid' :: 'Tid'(),
                      'toTid' :: 'Tid'(),
                      'type' :: string() | binary(),
                      'extraList' :: list(),
                      'error' :: 'TimError'(),
                      'extraMap' :: dict:dict()}).
-type 'TimIqBean'() :: #'TimIqBean'{}.

%% struct 'TimRoster'

-record('TimRoster', {'subscription' :: string() | binary(),
                      'tid' = #'Tid'{} :: 'Tid'(),
                      'name' :: string() | binary(),
                      'extraMap' :: dict:dict()}).
-type 'TimRoster'() :: #'TimRoster'{}.

%% struct 'TimRemoteUserBean'

-record('TimRemoteUserBean', {'error' :: 'TimError'(),
                              'ub' :: 'TimUserBean'(),
                              'extraMap' :: dict:dict()}).
-type 'TimRemoteUserBean'() :: #'TimRemoteUserBean'{}.

%% struct 'TimRemoteRoom'

-record('TimRemoteRoom', {'error' :: 'TimError'(),
                          'room' :: 'TimRoom'(),
                          'extraMap' :: dict:dict()}).
-type 'TimRemoteRoom'() :: #'TimRemoteRoom'{}.

%% struct 'TimResponseBean'

-record('TimResponseBean', {'threadId' :: string() | binary(),
                            'error' :: 'TimError'(),
                            'extraList' :: list(),
                            'extraMap' :: dict:dict()}).
-type 'TimResponseBean'() :: #'TimResponseBean'{}.

%% struct 'TimSock5Bean'

-record('TimSock5Bean', {'fromTid' = #'Tid'{} :: 'Tid'(),
                         'toTid' = #'Tid'{} :: 'Tid'(),
                         'addr' :: string() | binary(),
                         'port' :: integer(),
                         'transport' :: integer(),
                         'pubId' :: string() | binary(),
                         'extraMap' :: dict:dict()}).
-type 'TimSock5Bean'() :: #'TimSock5Bean'{}.

%% struct 'TimSock5Bytes'

-record('TimSock5Bytes', {'pubId' :: string() | binary(),
                          'index' :: integer(),
                          'bytes' = [] :: list(),
                          'extraMap' :: dict:dict()}).
-type 'TimSock5Bytes'() :: #'TimSock5Bytes'{}.

%% struct 'TimPage'

-record('TimPage', {'fromTimeStamp' :: string() | binary(),
                    'toTimeStamp' :: string() | binary(),
                    'limitCount' :: integer(),
                    'extraMap' :: dict:dict()}).
-type 'TimPage'() :: #'TimPage'{}.

%% struct 'TimMessageIq'

-record('TimMessageIq', {'tidlist' :: list(),
                         'timPage' :: 'TimPage'(),
                         'midlist' :: list(),
                         'extraMap' :: dict:dict()}).
-type 'TimMessageIq'() :: #'TimMessageIq'{}.

%% struct 'TimAuth'

-record('TimAuth', {'domain' :: string() | binary(),
                    'username' :: string() | binary(),
                    'pwd' :: string() | binary()}).
-type 'TimAuth'() :: #'TimAuth'{}.

%% struct 'TimMBeanList'

-record('TimMBeanList', {'threadId' :: string() | binary(),
                         'timMBeanList' :: list(),
                         'reqType' :: string() | binary(),
                         'extraMap' :: dict:dict()}).
-type 'TimMBeanList'() :: #'TimMBeanList'{}.

%% struct 'TimPBeanList'

-record('TimPBeanList', {'threadId' :: string() | binary(),
                         'timPBeanList' :: list(),
                         'reqType' :: string() | binary(),
                         'extraMap' :: dict:dict()}).
-type 'TimPBeanList'() :: #'TimPBeanList'{}.

%% struct 'TimPropertyBean'

-record('TimPropertyBean', {'threadId' :: string() | binary(),
                            'interflow' :: string() | binary(),
                            'tls' :: string() | binary()}).
-type 'TimPropertyBean'() :: #'TimPropertyBean'{}.

-endif.
