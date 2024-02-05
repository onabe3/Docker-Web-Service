import { Inter } from 'next/font/google'
import { useState, useEffect } from 'react'

const inter = Inter({ subsets: ['latin'] })

interface Result {
  message?: string,
  error?: string,
}

export default function Home() {
  const [containerName, setContainerName] =useState<string>("")
  const [os, setOs] =useState("");
  const [result, setResult] = useState<Result | undefined>(undefined);
  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
 
    // const formData = new FormData(event.currentTarget)
    const res = await fetch('/api/build', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({createOs: os, createName: containerName})
    })
    const data = await res.json()
    setResult(data)
  }

  useEffect(() => {
    console.log(result);
    
  }, [result])
  return (
    <>
      <form onSubmit={onSubmit}>
        <div className='flex flex-col'>
          <div>
            <label>
              <input name="os" type="radio" onChange={() => setOs("centos")} />
                Centos
            </label>
            <label>
              <input name="os" type="radio" onChange={() => setOs("ubuntu")} />
                Ubuntu
            </label>
            <label>
              <input name="os" type="radio" onChange={() => setOs("debian")} />
                Debian
            </label>
            <label>
              <input name="os" type="radio" onChange={() => setOs("almalinux")} />
                almalinux
            </label>
          </div>
          <input type="text" name="container_name" onChange={(e) => setContainerName(e.target.value)} />
          <button type='submit'>submit</button>
        </div>
      </form>
      {result === undefined ? "" : result?.message.length != 0 ? <p>成功しました。コンテナIDは{result?.message}</p> : <span color="red">失敗しました。{result?.error}</span>}
    </>
  )
}
