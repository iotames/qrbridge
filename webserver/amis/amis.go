package amis

type JsonContent any

type BaseComponent struct {
	Type string `json:"type"`
}

// {
//           type: 'page',
//           title: '客户PO文件格式转换',
//           body: {
//             type: 'form',
//             mode: 'horizontal',
//             api: '/api/potransform',
//             body: [
//               {
//                 "label": "客户简称",
//                 "type": "select",
//                 "name": "inputtpl",
//                 "source": "/api/customer/list"
//               },
//               {
//                 "type": "input-file",
//                 "name": "inputfile",
//                 "accept": ".xlsx",
//                 "label": "上传.xlsx文件",
//                 "maxSize": 10048576,
//                 "receiver": "/api/uploadfile"
//               }
//             ]
//           }
// }
