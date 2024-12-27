export function load({ cookies }) {
  const useruuid = cookies.get('useruuid')
  const visited = cookies.get('visited')

  console.log('cookies.set')
  cookies.set('visited', 'true', { path: '/' })

  return {
    visited: visited === 'true',
    useruuid: useruuid
  }
}
