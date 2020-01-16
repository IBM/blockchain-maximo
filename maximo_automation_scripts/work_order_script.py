from psdi.util.logging import FixedLoggers
from org.apache.http.impl.client import DefaultHttpClient
from org.apache.http import HttpEntity, HttpHeaders, HttpResponse, HttpVersion
from org.apache.http.entity import StringEntity
from org.apache.http.client.methods import HttpPost
from psdi.iface.router import HTTPHandler
from java.util import HashMap
from java.lang import String
from com.ibm.json.java import JSONObject
from com.ibm.json.java import JSONArray

FixedLoggers.MAXIMOLOGGER.info("####starting workorder script")

asset = mbo.getAsset()
asset_num = asset.getString("ASSETNUM")

mboSet = mbo.getThisMboSet()

wonum = mbo.getString("wonum") #mbo.getString("WOObjectName")
FixedLoggers.MAXIMOLOGGER.info("wonum")
FixedLoggers.MAXIMOLOGGER.info(wonum)

wo_asset_num = mboSet.getString("ASSETNUM")
FixedLoggers.MAXIMOLOGGER.info(wo_asset_num)

vendor = mboSet.getString("VENDOR")
FixedLoggers.MAXIMOLOGGER.info(vendor)

wo_status = mboSet.getString("STATUS")
FixedLoggers.MAXIMOLOGGER.info(wo_status)

args = JSONArray()
wo_num = "1"
args.add("workorder" + wonum)
args.add(wo_status)
args.add(vendor)
# add asset reference if included. TODO, add functionality to register asset props in blockchain
if asset_num:
    args.add("asset" + asset_num)

ctorMsg = JSONObject()
ctorMsg.put("function", "init_work_order")
ctorMsg.put("args", args)
params = JSONObject()
params.put("ctorMsg", ctorMsg)
obj = JSONObject()
obj.put('method', 'invoke')
obj.put('params', params)
jsonStr = obj.serialize(True)

FixedLoggers.MAXIMOLOGGER.info("json obj")
FixedLoggers.MAXIMOLOGGER.info(obj)

# post json to chaincode
handler = HTTPHandler()
map = HashMap()
url = "http://c1dfe6f7.ngrok.io"
map.put("URL", url + "/api/chaincode")
map.put("HTTPMETHOD", "POST")
map.put("body", jsonStr)
map.put("headers", "Content-Type: application/json")


# init HTTP Client, post JSON to blockchain server
client = DefaultHttpClient()
request = HttpPost(url + "/api/chaincode")
request.addHeader(HttpHeaders.CONTENT_TYPE, "application/json")
request.addHeader(HttpHeaders.ACCEPT, "application/json")
entity = StringEntity(jsonStr, "UTF-8")
request.setEntity(entity)
response = client.execute(request)
status = response.getStatusLine().getStatusCode()


# location = mboSet.getString("LOCATION")
# FixedLoggers.MAXIMOLOGGER.info(location)

# site = mbo.getString("SITE")
# FixedLoggers.MAXIMOLOGGER.info(site)

# using WO 1393
