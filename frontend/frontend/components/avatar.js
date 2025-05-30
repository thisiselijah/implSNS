export default function Avatar( props ) {
  let username = props.username ? props.username: "Default User";
  return (
    <>
      <div className="bg-white p-5 rounded-lg shadow">
        <h3 className="text-lg font-semibold text-gray-700 mb-3">{username}</h3>
        <div className="flex -space-x-2 overflow-hidden">
        <img
          alt=""
          src="https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
          className="inline-block size-10 rounded-full ring-2 ring-white"
        />
      </div>
        
      </div>
      
    </>
  )
}
