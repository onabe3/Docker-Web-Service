// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'

type Data = {
  name: string
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>
) {
  if(req.method==='POST') {
    const os = req.body.createOs
    const name = req.body.createName
    
    const result = await fetch(`http://127.0.0.1:1323/create/${os}/${name}`, {
      method: "POST",
    })
      .then(res => {
        return res.json();
      })
      .then(data => {
        return data;
      })
    res.status(200).json(result)
  }
}
